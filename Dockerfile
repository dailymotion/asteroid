# Alpine is used to be able to go inside the container for debug purposes
FROM golang:alpine

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Tini is used to better handles signals from the applications and have a better host stability
RUN apk add --no-cache tini

# Create group, user and folders
# UID and GID have been put above 10000 for security reasons
RUN addgroup --gid 11000 -S asteroid && adduser --uid 11000 -S asteroid -G asteroid  \
    && mkdir -p /home/asteroid/build \
    && mkdir /home/asteroid/dist \
    && chown -R asteroid:asteroid /home/asteroid

# Tell docker that all future commands should run as the asteroid user
USER asteroid

# Move to build directory
WORKDIR /home/asteroid/build

# Copy the code into the container build folder
COPY . .

# Copy config to asteroid user $HOME
COPY pkg/config/asteroid_example.yaml /home/asteroid/.asteroid.yaml

# Download dependencies using go mod
RUN go mod download

# Build the application with specific ENV
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -o asteroid ./cmd/asteroid

# Move to dist directory 
WORKDIR /home/asteroid/dist

# Copy the binary to dist folder
RUN cp /home/asteroid/build/asteroid . 

# ENTRYPOINT allow us to run the executable and pass arguments at run time.
ENTRYPOINT ["/sbin/tini", "--", "./asteroid"]
