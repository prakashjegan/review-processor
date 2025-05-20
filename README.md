# Review Processor

A robust S3-based review file processing system that processes hotel reviews from S3 files and maintains aggregated metrics.

## Project Overview

This system is designed to process hotel reviews stored in S3 files. It follows a microservices architecture pattern and uses event-driven processing to handle review data efficiently. The system reads review files from S3, processes them through a transformation pipeline, and maintains aggregated metrics for hotels.

## System Flow

1. **Scheduler**
   - Runs daily to fetch new review files from S3
   - Triggers the review processing pipeline

2. **Provider Factory**
   - Creates appropriate file handlers based on the file format
   - Implements factory pattern for extensibility

3. **File Processing**
   - Each file format has its specific implementation
   - Handles file-specific data formats and transformations

4. **S3 Integration**
   - Reviews are stored in AWS S3
   - System reads review files from configured S3 buckets
   - Maintains file processing state and checksums

5. **Redis Stream**
   - Acts as a message queue for review processing
   - Ensures reliable processing of reviews
   - Handles backpressure and retries

6. **Stream Consumer**
   - Consumes messages from Redis stream
   - Processes reviews through the transformation pipeline

7. **Data Persistence**
   - Transforms and stores review data in PostgreSQL
   - Maintains aggregated metrics for hotels

## Project Structure

```
.
├── app/
│   ├── config/                 # Configuration management
│   ├── database/              # Database connection and setup
│   ├── review/                # Core review processing logic
│   │   └── main.go           # Main application entry point
│   └── review-services/       # Review service implementations
│       ├── aws/              # AWS integration
│       ├── database/         # Database models and DAOs
│       ├── models/           # Business objects
│       ├── redisstream/      # Redis stream processing
│       └── scheduler/        # Job scheduling
├── config/                    # Configuration files
└── pkg/                      # Shared packages
```

## Data Models

### 1. ReviewRaw
- Raw review data from S3 files
- Contains unprocessed review information
- Tracks processing status and metadata

### 2. ProductReview
- Aggregated review data for a hotel
- Contains overall scores and metrics
- Maintains review counts

### 3. ReviewComment
- Individual review comments
- Stores detailed review information
- Links to reviewer information

### 4. ReviewerInfo
- Reviewer details and metadata
- Contains reviewer information

### 5. ProductReviewByProvider
- Provider-specific review metrics
- Maintains scores and counts per provider

### 6. ProductReviewGrade
- Detailed grading information
- Stores specific aspects and ratings

## Prerequisites

1. **Go Environment**
   - Go 1.16 or higher
   - GOPATH properly configured

2. **Database**
   - PostgreSQL 12 or higher
   - Required extensions: gorm

3. **Redis**
   - Redis 6.0 or higher
   - Redis Stream support enabled

4. **AWS**
   - AWS Account with appropriate permissions
   - S3 bucket for review storage
   - AWS credentials configured

5. **Other Dependencies**
   - GORM for database operations
   - Redis client for Go
   - AWS SDK for Go

## Setup and Execution

1. **Configuration**
   ```bash
   # Set up environment variables
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_NAME=review_db
   export DB_USER=postgres
   export DB_PASSWORD=your_password
   
   export REDIS_HOST=localhost
   export REDIS_PORT=6379
   
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_REGION=your_region
   export S3_BUCKET_NAME=your-review-bucket

   # Optional: Skip table creation if tables already exist
   export CREATE_TABLE_IGNORE=true
   ```

2. **Database Setup**
   ```bash
   # Create database
   createdb review_db
   ```

3. **Build and Run**
   ```bash
   # Build the project
   go build -o review-processor

   # Run the service
   ./review-processor
   ```

   The main application (`app/review/main.go`) will:
   - Initialize the database connection
   - Create required tables if they don't exist (unless CREATE_TABLE_IGNORE=true)
   - Start the review processing service

## Notes

1. **Unit Tests**
   - Unit tests are not implemented yet
   - Test coverage needs to be added

2. **Docker Support**
   - Docker configuration is not implemented
   - Dockerfile and docker-compose.yml need to be created

3. **Monitoring**
   - Basic logging is implemented
   - Metrics collection needs to be added

4. **Error Handling**
   - Basic error handling is in place
   - More comprehensive error handling needed

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License
