# Build stage
FROM node:18.17.0 as build

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies and set up environment
RUN npm ci && \
    chmod -R 755 /app/node_modules/.bin && \
    export PATH="/app/node_modules/.bin:$PATH"

# Copy source code
COPY . .

# Build the app
ENV NODE_ENV=production
ENV PATH="/app/node_modules/.bin:${PATH}"
SHELL ["/bin/bash", "-c"]
RUN ./node_modules/.bin/react-scripts build

# Production stage
FROM nginx:alpine

# Copy built files from build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-80}/ || exit 1

EXPOSE ${PORT:-80}

CMD ["nginx", "-g", "daemon off;"]
