const config = {
    apiVersion: process.env.APP_API_VERSION || "v1", 
    port: process.env.BACKEND_PORT || 5000, 
};

config.backendBaseUrl = `http://localhost:${config.port}/${config.apiVersion}`;

export default config;
