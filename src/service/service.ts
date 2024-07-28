// Service logic, singleton
class Service {
    private static _instance: Service;

    // Constructor
    private constructor() {}

    // Global access point
    public static getInstance(): Service {
        if (!Service._instance) {
            Service._instance = new Service();
        }
        return Service._instance;
    }

    // $ lucy init
    // public initialization(): void {}
}

// File related jobs
// Network related jobs: fetch api, download, etc.

export { Service };
