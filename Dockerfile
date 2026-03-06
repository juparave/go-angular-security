# ============================================================
# Stage 1: Build Go backend
# ============================================================
FROM golang:1.24-alpine AS build_back

RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/bin/server ./cmd/api/

# ============================================================
# Stage 2: Build Angular frontend
# ============================================================
FROM node:22-alpine AS frontend

WORKDIR /app/angular
COPY angular/package.json angular/package-lock.json ./
RUN npm ci

COPY angular/ ./
RUN npx ng build --configuration production

# ============================================================
# Stage 3: Final runtime image
# ============================================================
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy backend binary
COPY --from=build_back /app/bin/server ./server

# Copy frontend dist
COPY --from=frontend /app/angular/dist/angular ./dist/angular

# Copy email templates (embedded in binary, but keep for reference)
# Templates are embedded at build time via go:embed

# Run as unprivileged user
USER app

EXPOSE 5000

CMD ["./server"]
