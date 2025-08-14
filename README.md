# IBM&reg; Power&reg; Access Cloud

![pac-design](https://github.com/PDeXchange/pac/assets/4647967/9c2fc7a5-7006-4abe-9b3c-525439a738d7)


# PAC Dev Setup Guide
This guide provides step-by-step instructions for setting up and running the Power Access Cloud locally.

## Prerequisites

* Git
* Go
* Node.js and npm
* Yarn package manager

### Part 1: PAC Server Setup

1. Clone Repository

    ```
    git clone https://github.com/PDeXchange/pac.git
    cd pac
    ```

2. Set Up Kubernetes Cluster
    
    Bring up any Kubernetes cluster (minikube, kind, or cloud-based cluster).

3. Generate Manifests and Deploy

    ```
    make generate && make manifests
    make deploy
    ```

4. Configure Environment Variables

    ```
    export KEYCLOAK_CLIENT_ID=your_client_id
    export KEYCLOAK_CLIENT_SECRET=your_client_secret
    export KEYCLOAK_REALM=your_realm
    export KEYCLOAK_HOSTNAME=your_keycloak_hostname
    export KEYCLOAK_SERVICE_ACCOUNT=your_service_account
    export KEYCLOAK_SERVICE_ACCOUNT_PASSWORD=your_service_account_password
    export MONGODB_URI=your_mongodb_connection_string
    ```

5. Start the Backend Server
    
    ```
    go run main.go
    ```

### Part 2: PAC UI Setup

1. Clone Repository

    ```
    git clone https://github.com/PDeXchange/pac-ui.git
    cd pac
    ```

2. Configure Environment Variables

    Create a .env file or export the following variables:
    ```
    export REACT_APP_KEYCLOAK_URL=your_keycloak_url
    export REACT_APP_KEYCLOAK_REALM=your_keycloak_realm
    export REACT_APP_KEYCLOAK_CLIENT_ID=your_keycloak_client_id
    export REACT_APP_PAC_GO_SERVER_TARGET=your_backend_server_url
    ```

    and run `./env.sh`

3. Install Dependencies and Run

    ```
    yarn install
    npm start
    ```

    The UI will be available at http://localhost:3000

### Part 3: PAC Controller Setup
1. Create IBM Cloud API Key
2. Configure and Run Controller

    ```
    export IBMCLOUD_APIKEY=your_ibm_cloud_api_key
    cd pac
    make run
    ```