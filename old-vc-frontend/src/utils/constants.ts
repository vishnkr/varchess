export interface GameMode{
    name: string,
    key:string
}
export const GAME_MODES : GameMode[] = [
    {
        name: "Standard 8x8",
        key: "standard"
    },
    {
        name: "Custom Variant",
        key: "custom"
    },
    {
        name: "Duck Chess",
        key: "duckChess"
    },
]
