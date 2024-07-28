import { Server } from "./serverType.js";

// Service logic, singleton
class Service {
    private static service: Service;
    private server: Server = new Server();

    // Constructor
    private constructor() {
        this.server = new Server();
        this.detectServerModules();
    }

    // Global access point
    public static getService(): Service {
        if (!Service.service) {
            Service.service = new Service();
        }
        return Service.service;
    }

    // Assign values to server.modules
    private detectServerModules(): void {
        // Detect MCDR
        this.server.modules.MCDR = {};
        // Detect Minecraft
        this.server.modules.Minecraft = {};
        // Detect Fabric
        // Detect Forge
        // Detect Bukkit
    }

    // $ lucy init
    public initialization(): void {}
}

// File related jobs
// Network related jobs: fetch api, download, etc.

export { Service };
