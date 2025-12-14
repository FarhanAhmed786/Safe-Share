# Safe Share

Safe Share is a secure, temporary file sharing solution designed for quick, one-time file transfers. It allows users to upload sensitive files which are encrypted and stored in the server's volatile memory. A unique, single-use link is generated for the recipient. Once the file is downloaded, it is permanently deleted from the server.

## Usage

1.  **Upload**: Navigate to the home page, select a file (up to 10MB), and click "Upload".
2.  **Share**: Copy the generated link provided after a successful upload.
3.  **Download**: Send the link to the intended recipient. When they access the link, the file is decrypted and downloaded to their device.
4.  **Vanish**: Immediately after the download starts, the file is removed from the server's memory. The link becomes invalid.

## Getting Started

### Prerequisites

-   Go 1.25 or higher
-   OpenSSL (for generating self-signed certificates)
-   Make

### Running the Application

The project includes a `Makefile` to handle SSL certificate generation and application startup.

1.  **Start the Server**:
    Run the following command in the project root. This will generate self-signed SSL certificates (`server.crt`, `server.key`) if they do not exist and start the Go server.
    ```bash
    make run
    ```

2.  **Access the Application**:
    Open your web browser and navigate to:
    `https://localhost:8080`

    *Note: Since self-signed certificates are used, your browser may display a security warning. You will need to accept the risk to proceed to localhost.*

3.  **Clean Up**:
    To remove the generated SSL certificates:
    ```bash
    make clean
    ```

## Architecture

The application is built as a monolithic HTTP server using Go's standard library.

1.  **Upload Handling**: Incoming files are read into memory. To prevent denial-of-service attacks, request bodies are limited to 10MB.
2.  **Encryption**: Before storage, files are encrypted using AES-256 in GCM mode. A unique 32-byte key and 12-byte nonce are generated for each file.
3.  **Storage**: The encrypted binary data, along with its decryption key and metadata, is stored in a thread-safe in-memory map (`FileStore`). The server enforces a global RAM usage limit of 50MB.
4.  **Token Generation**: A cryptographically secure random token is generated and acts as the key to retrieve the file.
5.  **Retrieval & Decryption**: When a valid token is presented via the download endpoint, the file is retrieved, decrypted on-the-fly, and streamed to the client.
6.  **Cleanup**: The file entry is immediately deleted from the memory map upon retrieval, freeing up server resources.

## Limitations

-   **Volatile Storage**: Files are stored exclusively in the server's RAM. If the application restarts or crashes, all uploaded files are lost.
-   **Memory Constraints**: The server has a hard limit on total memory usage (50MB). Heavy usage may result in rejected uploads until space is freed.
-   **Link Security**: The security model relies entirely on the secrecy of the generated link. If an attacker intercepts the link before the intended recipient, they can download the file. Since the link is one-time use, the intended recipient would then find the link invalid, alerting them to the breach.
-   **Scalability**: Due to in-memory storage, this solution is not suitable for large file transfers or high-concurrency environments without vertical scaling of RAM.

## Future Scope

-   **Persistent Storage**: Integration with a database (SQL or NoSQL) or object storage (S3) to allow for larger file retention and persistence across server restarts.
-   **Expiration Time**: Implementing a time-based expiration (TTL) for files that are not downloaded within a specific window.
-   **Password Protection**: Adding an optional password field during upload to require both the token link and a password for decryption.

## Developer Note

This project is a Proof of Concept (PoC) developed to explore the mechanics of secure, ephemeral file sharing. It is intended for educational purposes and architectural demonstration.

For a production-grade deployment, the following enhancements are necessary:
-   **Persistence**: Transitioning from in-memory storage to a robust database system to handle data reliability and scale.
-   **Infrastructure**: Deployment on a secured public server rather than localhost.
-   **Authentication**: Implementing stricter access controls, such as requiring a secondary authentication factor (e.g., a password shared separately, or an authenticator app).