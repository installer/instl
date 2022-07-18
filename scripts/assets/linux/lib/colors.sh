# Foreground Colors

fRed () {
    printf "\e[31m%s\e[0m" "$1"
}

fRedLight () {
    printf "\e[91m%s\e[0m" "$1"
}

fYellow () {
    printf "\e[33m%s\e[0m" "$1"
}

fYellowLight () {
    printf "\e[93m%s\e[0m" "$1"
}

fGreen () {
    printf "\e[32m%s\e[0m" "$1"
}

fGreenLight () {
    printf "\e[92m%s\e[0m" "$1"
}

fBlue () {
    printf "\e[34m%s\e[0m" "$1"
}

fBlueLight () {
    printf "\e[94m%s\e[0m" "$1"
}

fMagenta () {
    printf "\e[35m%s\e[0m" "$1"
}

fMagentaLight () {
    printf "\e[95m%s\e[0m" "$1"
}

fCyan () {
    printf "\e[36m%s\e[0m" "$1"
}

fCyanLight () {
    printf "\e[96m%s\e[0m" "$1"
}

fWhite () {
    printf "\e[37m%s\e[0m" "$1"
}

## Background Colors
bRed () {
    printf "\e[41m%s\e[0m" "$1"
}

bRedLight () {
    printf "\e[101m%s\e[0m" "$1"
}

bYellow () {
    printf "\e[43m%s\e[0m" "$1"
}

bYellowLight () {
    printf "\e[103m%s\e[0m" "$1"
}

bGreen () {
    printf "\e[42m%s\e[0m" "$1"
}

bGreenLight () {
    printf "\e[102m%s\e[0m" "$1"
}

bBlue () {
    printf "\e[44m%s\e[0m" "$1"
}

bBlueLight () {
    printf "\e[104m%s\e[0m" "$1"
}

bMagenta () {
    printf "\e[45m%s\e[0m" "$1"
}

bMagentaLight () {
    printf "\e[105m%s\e[0m" "$1"
}

bCyan () {
    printf "\e[46m%s\e[0m" "$1"
}

bCyanLight () {
    printf "\e[106m%s\e[0m" "$1"
}

bWhite () {
    printf "\e[47m%s\e[0m" "$1"
}

# Special Colors

resetColor () {
    printf "\e[0m"
}

# Theme

primaryColor () {
   fCyan "$1"
}

secondaryColor () {
   fMagentaLight "$1"
}

info () {
   fBlueLight " i " && resetColor && fBlue "$1" && echo
}
