import * as fs from "fs";

export async function initialization() {
    fs.mkdirSync(".mcpm");
}
