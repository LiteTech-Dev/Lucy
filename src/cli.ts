import { Command } from "commander";
import inquirer from "inquirer";
import { prompts } from "./inquirer-prompts.js";

const program = new Command();

program
    .name("mcpm")
    .description("")
    .action(() => {
        console.log(
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
    .description("")
    .action(() => {});

program.command("list");

program.command("install");

program.command("remove");

program.command("update");

program.command("disable");

program.command("enable");

program.command("export");

program.command("import");

program.command("host");

export { program };
