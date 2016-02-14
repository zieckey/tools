syntax on

let $LANG="en_US.UTF-8"
set fileencodings=utf-8,chinese,latin-1
set termencoding=utf-8
set encoding=utf-8

set tabstop=4
set softtabstop=4
set autoindent
set hlsearch "hight light search word
set incsearch
set cindent
set shiftwidth=4 "auto-indent amount when using cindent, >>, << and stuff like that"
set cinoptions={0,1s,t0,n-2,p2s,(03s,=.5s,>1s,=1s,:1s
set foldmethod=marker
if &term=="xterm"
    set t_Co=8
    set t_Sb=^[[4%dm
    set t_Sf=^[[3%dm
endif

"set columns=80
set shiftwidth=4  
set expandtab


set nocp
filetype plugin on

"using 4 space to substitute TAB
autocmd FileType c,cpp,py,h,cxx,hxx,hpp,cc,php,CC,C set shiftwidth=4 | set expandtab

au! BufRead,BufNewFile *.thrift setfiletype thrift
au! BufRead,BufNewFile *.json setfiletype json
au! BufRead,BufNewFile *.proto setfiletype proto
au! BufRead,BufNewFile *.go setfiletype go
"autocmd FileType go autocmd BufWritePre <buffer> Fmt


set fileformats=unix,mac,dos
"set fileencodings=utf8,gb1803,gbk,cp936,iso-8859-1

call pathogen#infect()

syntax enable
filetype plugin on
"set number
let g:go_disable_autoinstall = 0
