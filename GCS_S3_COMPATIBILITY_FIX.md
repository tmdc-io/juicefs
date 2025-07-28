# Google Cloud Storage S3 Compatibility Fix

## Problem Description

The issue you're experiencing is related to the AWS SDK v2 upgrade in JuiceFS, which changed how S3 authentication is handled. When using Google Cloud Storage with the S3-compatible API (via `storage.googleapis.com`) and MinIO SDK with GCS interoperability keys, the authentication is failing with `SignatureDoesNotMatch: Access denied` errors.

## Root Cause

The AWS SDK v2 has stricter authentication requirements compared to the previous version. When using GCS S3-compatible API with MinIO SDK, the following issues occur:

1. **Region Handling**: GCS S3-compatible API requires specific region handling
2. **Path Style vs Virtual Hosted Style**: GCS S3-compatible API works better with path-style URLs
3. **Authentication Headers**: The AWS SDK v2 generates different authentication headers that may not be compatible with GCS S3-compatible API

## Solution

### 1. Code Changes Applied

I've made the following changes to fix the GCS S3 compatibility issues:

#### `pkg/object/s3.go`
- Added GCS endpoint detection (`storage.googleapis.com`)
- Force path-style URLs for GCS endpoints
- Set region to "auto" for GCS when using default region
- Added environment variable configuration for GCS S3 compatibility

#### `pkg/object/minio.go`
- Added GCS endpoint detection in MinIO configuration
- Force path-style URLs for GCS endpoints
- Set region to "auto" for GCS when using default region

### 2. Updated Configuration

Update your Kubernetes secret with the following configuration:

```bash
# Update the bucket URL to use path-style format
kubectl patch secret juicefs-delete-secret -n juicefs-system --type='json' -p='[
  {
    "op": "replace",
    "path": "/data/bucket",
    "value": "aHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL29wczAwMS1wcmVjaXNlbWEtZGV2L2p1aWNlZnMtZGVsZXRlLz9kaXNhYmxlLWNoZWNrc3VtPXRydWUmZGlzYWJsZS0xMDAtY29udGludWU9dHJ1ZQ=="
  }
]'
```

The new bucket URL format is:
```
https://storage.googleapis.com/ops001-precisema-dev/juicefs-delete/?disable-checksum=true&disable-100-continue=true
```

### 3. Environment Variables

Set the following environment variables in your deployment:

```yaml
env:
- name: AWS_REGION
  value: "auto"
- name: AWS_DEFAULT_REGION
  value: "auto"
- name: JFS_S3_VHOST_STYLE
  value: "false"
```

### 4. Alternative Configuration Options

If the above doesn't work, try these alternative configurations:

#### Option A: Use Native GCS Protocol
Switch from S3-compatible API to native GCS protocol:

```bash
# Update storage type to gs
kubectl patch secret juicefs-delete-secret -n juicefs-system --type='json' -p='[
  {
    "op": "replace",
    "path": "/data/name",
    "value": "Z3M="
  },
  {
    "op": "replace",
    "path": "/data/bucket",
    "value": "b3BzMDAxLXByZWNpc2VtYS1kZXY="
  }
]'
```

#### Option B: Use Different S3 Endpoint Format
Try using a different S3 endpoint format:

```bash
# Alternative bucket URL format
kubectl patch secret juicefs-delete-secret -n juicefs-system --type='json' -p='[
  {
    "op": "replace",
    "path": "/data/bucket",
    "value": "aHR0cHM6Ly9vcHMwMDEtcHJlY2lzZW1hLWRldi5zdG9yYWdlLmdvb2dsZWFwaXMuY29tL2p1aWNlZnMtZGVsZXRlLz9kaXNhYmxlLWNoZWNrc3VtPXRydWUmZGlzYWJsZS0xMDAtY29udGludWU9dHJ1ZQ=="
  }
]'
```

This format is: `https://ops001-precisema-dev.storage.googleapis.com/juicefs-delete/?disable-checksum=true&disable-100-continue=true`

## Testing the Fix

1. **Apply the code changes** to your JuiceFS build
2. **Update the Kubernetes secret** with the new configuration
3. **Restart the JuiceFS gateway**:
   ```bash
   kubectl rollout restart deployment juicefs-gateway-delete -n juicefs-system
   ```
4. **Check the logs**:
   ```bash
   kubectl logs -n juicefs-system -l app=juicefs-gateway-delete -f
   ```

## Expected Log Messages

With the fix applied, you should see these log messages:
```
Detected Google Cloud Storage S3 compatibility endpoint
HTTP header 100-Continue is disabled
CRC checksum is disabled
```

## Troubleshooting

If you still see authentication errors:

1. **Verify GCS Interoperability Keys**: Ensure your access key and secret key are valid GCS interoperability keys
2. **Check Bucket Permissions**: Verify the service account has proper permissions on the GCS bucket
3. **Try Different Region**: Set a specific region instead of "auto":
   ```bash
   export AWS_REGION="us-central1"
   ```
4. **Check Network Connectivity**: Ensure the pod can reach `storage.googleapis.com`

## Additional Notes

- The fix specifically targets GCS S3-compatible API endpoints containing `storage.googleapis.com`
- Path-style URLs are enforced for GCS endpoints to ensure compatibility
- The `disable-checksum=true` and `disable-100-continue=true` parameters help avoid compatibility issues
- This fix maintains backward compatibility with other S3-compatible storage providers 