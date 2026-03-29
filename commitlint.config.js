export default {
    extends: ["@commitlint/config-conventional"],
    rules: {
        "breaking-change-forbidden": [2, "always"],
        "type-enum": [2, "never", ["!"]], // doesn't work for the ! syntax
    },
    plugins: [
        {
            rules: {
                "breaking-change-forbidden": ({
                    body,
                    footer,
                    subject,
                    type,
                }) => {
                    const hasBreakingFooter = /BREAKING[\-\ ]CHANGE/i.test(
                        `${body ?? ""} ${footer ?? ""}`,
                    );
                    const hasBreakingType =
                        type?.endsWith("!") || subject?.startsWith("!");

                    return [
                        !(hasBreakingFooter || hasBreakingType),
                        "Breaking changes are not allowed on this branch",
                    ];
                },
            },
        },
    ],
};
