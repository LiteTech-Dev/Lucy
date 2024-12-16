import { Command } from "commander";
import inquirer from "inquirer";
import { prompts } from "./inquirer-prompts.js";
import * as cmdInfo from "./commands/info/index.js";
import * as cmdInstall from "./commands/install/index.js";

const program = new Command();

program
    .name("lucy")
    .description("")
    .action(() => {
        console.warn(
            "A Lucy installation was not detected in the current directory."
        );
        inquirer.prompt(prompts.initConfirmation).then((answers) => {
            if (answers.confirm) {
                console.log("Placeholder for initialization");
            }
        });
    });

program
    .command("init")
    .description("Initialize Lucy in the current directory.")
    .action(() => {});

program
    .command("info")
    .description("Information on the current server.")
    .action(() => {
        cmdInfo;
        // TODO: Change this
    });

program
    .command("install")
    .description("Install mods, plugins, mod loaders, or even Minecraft.")
    .action(() => {
        cmdInstall;
    });

program
    .command("remove")
    .description("Safely remove mods, plugins, or mod loaders installation.");

program
    .command("update")
    .description(
        "Update the mods, plugins, or mod loaders to the latest version."
    );

program
    .command("migrate")
    .description("Attempt to migrate the server to another Minceraft version.");

program
    .command("disable")
    .description("Disable mods, plugins, or mod loaders.");

program
    .command("enable")
    .description("Enable mods, plugins, or mod loaders disabled before.");

// program
//     .command("export")
//     .description(
//         "Bundle the whole server so that it can be setup by Lucy on somewhere else."
//     );

// program.command("import").description("Import a server bundle.");

export { program };
