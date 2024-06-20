import * as fs from "fs";

export async function initialization(): Promise<void> {
    fs.mkdirSync(".mcpm");
    // const MCDR = serverProbe.checkMCDR;
}
