# Signing Service - Go Implementation

RESTful API for cryptographic signature devices with strictly monotonic counters and signature chaining (KassenSichV/RKSV compliance).

## Quick Start

Before running the server, make sure you have Go installed. Go version should be 1.20 or higher.

Check your Go version by running:

```bash
go version
```

If you don't have Go installed, you can download it from [here](https://golang.org/dl/).

Then you need to tidy up the dependencies by running:

```bash
go mod tidy
```

To run the server, use the following command:

```bash
go run main.go                              # Server on :8080
```

To run the tests, use the following command:

```bash
go test ./...                               # Run tests
go test -race ./...                         # Race detector
go test -coverprofile=coverage.out ./...    # Coverage report
```

## API

```bash
# Create device (ECC or RSA)
curl -X POST http://localhost:8080/api/v0/devices \
  -d '{"id":"device-1","algorithm":"ECC","label":"Register 1"}'

# Sign transaction
curl -X POST http://localhost:8080/api/v0/devices/device-1/sign \
  -d '{"data":"SALE:100.00:EUR"}'
# Returns: {"signature":"...", "signedData":"0_SALE:100.00:EUR_base64(deviceId)"}

# Get device
curl http://localhost:8080/api/v0/devices/device-1

# List all devices
curl http://localhost:8080/api/v0/devices
```

## Concurrency: Monotonic Counter

**Challenge:** Multiple concurrent clients → race conditions, counter gaps, invalid signatures.

**Solution:** Execute Around pattern with mutex-protected atomic updates.

```go
// ❌ Naive - race condition
device := repo.GetByID(id)
device.Counter++
repo.Update(device)

// ✅ Atomic update
repo.Update(deviceID, func(device *Device) error {
    securedData := fmt.Sprintf("%d_%s_%s", device.Counter, data, device.LastSignature)
    signature := sign(securedData)
    device.Counter++
    device.LastSignature = signature
    return nil
})
```

**Implementation:**

```go
func (r *InMemoryRepository) Update(
    deviceID string,
    updateFn func(*domain.Device) error,
) (*domain.Device, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    device := r.devices[deviceID]
    if err := updateFn(device); err != nil {
        return nil, err
    }
    return device, nil
}
```

**Why it works:** Mutex serializes all counter updates. Only one goroutine can read-sign-increment at a time.

## Testing

**Coverage:** Persistence 96.3% | Service 87.5% | API 57.1%

## Architecture

**Layers:** API → Service → Repository → Domain

**Key Patterns:**

- **Repository:** Interface-based data access, easy to swap in-memory for database
- **Dependency Injection:** Services depend on interfaces, improves testability
- **Execute Around:** Atomic updates with automatic mutex management
- **Factory:** Crypto algorithm selection (RSA/ECC)

**Design Decisions:**

- Single service layer (signing is a device operation)
- gorilla/mux for path parameters and HTTP method routing
- base64.RawStdEncoding for binary-safe JSON transmission

## AI Tools Usage

**Tool:** Cascade (Claude AI)

**Usage Breakdown:**

- **Architecture (20%):** Brainstorming layered design, repository pattern
- **Concurrency (15%):** Discussing Execute Around pattern for handling race condition
- **Code (5%):** Repetitive test boilerplate only
- **Debugging (25%):** Spotted bugs
- **Documentation (70%):** README structure and content

**Design Decisions:**

- **Execute Around pattern** → ensures atomic read-modify-write
  - Trade-off: Less flexible than separate lock/unlock, but safer (impossible to forget unlock)
- **Repository pattern** → easy database migration later
  - Trade-off: Extra abstraction layer, but decouples business logic from storage
- **gorilla/mux** → path parameters and method routing out of the box
  - Trade-off: External dependency, but saves manual parsing and reduces boilerplate
- **Mutex serialization** → guarantees no counter gaps
  - Trade-off: Limits throughput (one sign at a time per device), but ensures correctness

## Assumptions & Limitations

**Assumptions:**

- Device IDs are unique and provided by client
- Single server instance (no distributed locking needed)

**Known Limitations:**

- In-memory storage (data lost on restart)
- No authentication/authorization
- No signature verification endpoint
- Mutex limits throughput to sequential signing per device

**Time Spent:** ~10 hours

---

**Author:** Husnul Nawafil  
**Date:** October 26, 2025  
**Challenge:** fiskaly Signing Service Coding Challenge
