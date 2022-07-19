# Foreground Colors

function fRed
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkRed
}

function fRedLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Red
}

function fGreen
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkGreen
}

function fGreenLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Green
}

function fYellow
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkYellow
}

function fYellowLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Yellow
}

function fBlue
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkBlue
}

function fBlueLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Blue
}

function fMagenta
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkMagenta
}

function fMagentaLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Magenta
}

function fCyan
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkCyan
}

function fCyanLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Cyan
}

function fWhite
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor White
}

function fBlack
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Black
}

function fGray
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor DarkGray
}

function fGrayLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -ForegroundColor Gray
}

# Background Colors

function bRed
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkRed
}

function bRedLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Red
}

function bGreen
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkGreen
}

function bGreenLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Green
}

function bYellow
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkYellow
}

function bYellowLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Yellow
}

function bBlue
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkBlue
}

function bBlueLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Blue
}

function bMagenta
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkMagenta
}

function bMagentaLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Magenta
}

function bCyan
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkCyan
}

function bCyanLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Cyan
}

function bWhite
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor White
}

function bBlack
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Black
}

function bGray
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor DarkGray
}

function bGrayLight
{
    param (
        $Msg
    )
    Write-Host $Msg -NoNewline -BackgroundColor Gray
}

# Special Colors

function resetColor
{
    Write-Host -NoNewline -Reset
}

# Theme

function primaryColor
{
    param (
        $Msg
    )
    fCyan $Msg
}

function secondaryColor
{
    param (
        $Msg
    )
    fMagentaLight $Msg
}

function info
{
    param (
        $Msg
    )
    fBlueLight " i "
    fBlue $Msg
    echo ""
}

function warning
{
    param (
        $Msg
    )
    fYellowLight " ! "
    fYellow $Msg
    echo ""
}

function error
{
    param (
        $Msg
    )
    fRedLight " X "
    fRed $Msg
    echo ""
}

function success
{
    param (
        $Msg
    )
    fGreenLight " + "
    fGreen $Msg
    echo ""
}

function verbose
{
    param (
        $Msg
    )
    if ($verbose)
    {
        fGrayLight " > "
        fGray $Msg
        echo ""
    }
}

