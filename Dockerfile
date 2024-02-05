# -------------
# build stage
# -------------
FROM golang:alpine AS build

# Workdir
WORKDIR /src

# Attach sources
ADD . /src

# System deps
RUN apk add --no-cache git npm

# Build
RUN go build -o news
RUN (cd static; npm i; npm run build)

# -------------
# runtime stage
# -------------
FROM alpine

# Workdir
WORKDIR /app

# Copy app
COPY --from=build /src/news /app/
COPY --from=build /src/*.html /app/
COPY --from=build /src/assets-dist /app/assets-dist

# Entrypoint
ENTRYPOINT ./news
