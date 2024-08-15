abstract class ServerError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "ServerError";
    }
}

abstract class ServerExecutableError extends ServerError {
    constructor(message: string) {
        super(message);
        this.name = "ServerExecutableError";
    }
}

class NoExecutableFoundError extends ServerExecutableError {
    constructor() {
        super("No server executable found");
        this.name = "NoExecutableFoundError";
    }
}

class MultipleExecutablesFoundError extends ServerExecutableError {
    constructor() {
        super("Multiple server executables found");
        this.name = "MultipleExecutablesFoundError";
    }
}

export {
    ServerExecutableError,
    NoExecutableFoundError,
    MultipleExecutablesFoundError,
};
