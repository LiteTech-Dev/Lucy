import * as fs from "fs";
import * as serverProbe from "./server-probe";

export async function initialization(): Promise<void> {
    const mcpmDirectory = ".mcpm";
    fs.mkdirSync(mcpmDirectory);
    // const MCDR = serverProbe.checkMCDR;
}
