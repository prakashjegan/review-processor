# S3 File Processor with Quartz Scheduler

This is a Go application that processes files from an S3 bucket using a Quartz scheduler. It maintains the state of processed files, handles retries, and tracks file processing events.

## Features

- S3 file listing and downloading
- Quartz scheduler for periodic file processing
- In-memory storage for file processing state
- File checksum calculation
- Retry mechanism for failed processing
- Event tracking for file processing status
- Configuration via environment variables

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── config/
│   └── config.go
├── internal/
│   ├── models/
│   │   └── file.go
│   ├── s3/
│   │   └── service.go
│   ├── scheduler/
│   │   └── scheduler.go
│   └── storage/
│       └── storage.go
├── .env
├── go.mod
└── README.md
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
# S3 Configuration
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=us-east-1
S3_BUCKET_NAME=your-bucket-name

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=file_processor

# Scheduler Configuration
SCHEDULER_INTERVAL_1=0 0/5 * * * ?  # Run every 5 minutes
MAX_RETRIES=3
RETRY_DELAY=300  # 5 minutes in seconds

# Application Configuration
LOG_LEVEL=info
ENVIRONMENT=development
```

## Building and Running

1. Install dependencies:
```bash
go mod download
```

2. Build the application:
```bash
go build -o file-processor ./cmd/server
```

3. Run the application:
```bash
./file-processor
```

## How it Works

1. The scheduler runs every 5 minutes (configurable)
2. It lists all files from the configured S3 bucket
3. For each file:
   - Checks if it's already processed
   - Downloads the file if not processed
   - Calculates checksum
   - Processes the file (custom logic can be added in the `processFile` method)
   - Updates the file status and events
4. Failed files are retried up to the configured maximum retries
5. All events and file statuses are tracked in memory

## Adding Custom Processing Logic

To add custom file processing logic, modify the `processFile` method in `internal/scheduler/scheduler.go`. This is where you can add your specific file processing requirements.

## Error Handling

- Failed file processing is tracked with error messages
- Files are retried up to the configured maximum retries
- All errors are logged
- File processing events are recorded for monitoring

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 