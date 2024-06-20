import { Command } from "commander";
import inquirer from "inquirer";

import { initialization } from "./initialization.js";
import { initConfirmation } from "./prompts.js";
import "./directory.js";
import { getMCPMDir } from "./directory.js";

const mcpm = new Command();

// Parent command
mcpm.name("mcpm")
    .description("Minecraft Pack Manager")
    .action(() => {
        console.log(getMCPMDir());
    }); // 本地化

// Sub-command `mcpm init`
mcpm.command("init")
    .description("Initialize MCPM at current directory")
    .action(async () => {
        const doInit = await inquirer.prompt(initConfirmation);
        if (doInit.confirm === true) {
            initialization();
            console.log("MCPM initialized.");
        } else {
            console.log("Stopped.");
        }
    });

// mcpm.command("install")
//     .description("Install a package")
//     .action(() => {});
// mcpm.command("update").description();
// mcpm.command("upgrade").description();

mcpm.parse(process.argv);
