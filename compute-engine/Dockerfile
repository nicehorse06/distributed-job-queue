FROM rust:1.85-alpine AS build

WORKDIR /app
RUN apk add --no-cache musl-dev

COPY proto ./proto
COPY compute-engine/Cargo.toml compute-engine/build.rs ./compute-engine/

WORKDIR /app/compute-engine

COPY compute-engine/src ./src
RUN cargo build --release

FROM alpine:3.20

RUN adduser -D -g '' appuser
USER appuser

COPY --from=build /app/compute-engine/target/release/compute-engine /bin/compute-engine

EXPOSE 50051

ENTRYPOINT ["/bin/compute-engine"]
