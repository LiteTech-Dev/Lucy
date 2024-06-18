import { Path } from "typescript";

declare global {
    let pwd: Path;
    let mcpmDirExists: boolean;
    let mcpmDirPath: Path | null;
    const mcpmDirName = ".mcpm";
}
