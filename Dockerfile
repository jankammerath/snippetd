# Use an ARM64 base image
FROM linuxkit/alpine:35b33c6b03c40e51046c3b053dd131a68a26c37a

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the host to the container
COPY ./bin/snippetd /app/snippetd

# Make the binary executable
RUN chmod +x /app/snippetd

# Expose the port the application listens on
EXPOSE 8504

# Command to run the application
CMD ["/app/snippetd"]