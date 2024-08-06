# Now create separate deployment image
FROM gcr.io/distroless/base

# Definition of this variable is used by 'skaffold debug' to identify a golang binary.
# Default behavior - a failure prints a stack trace for the current goroutine.
# See https://golang.org/pkg/runtime/
ENV GOTRACEBACK=single

WORKDIR /app
ADD ./echo-pipeline /app/echo-pipeline


# TODO: which port?
EXPOSE 80 

ENTRYPOINT ["/app/echo-pipeline"]
