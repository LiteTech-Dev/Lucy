import { Question } from "inquirer";

export const prompts = {
    initConfirmation: [
        {
            type: "confirm",
            name: "confirm",
            message: "Initialize at current directory?", // 本地化
            default: false,
        },
    ],
};
