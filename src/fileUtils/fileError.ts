abstract class FileError extends Error {
    constructor(message: string) {
        super(message);
        this.name = "FileError";
    }
}

class FileNotFoundError extends FileError {
    constructor() {
        super("File read error");
        this.name = "FileReadError";
    }
}

export { FileError, FileNotFoundError };
