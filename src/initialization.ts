import * as fs from "fs";
import * as serverProbe from "./server-probe";

export async function initialization(): Promise<void> {
    fs.mkdirSync(".mcpm");
    // const MCDR = serverProbe.checkMCDR;
}
