function send_email()
{
    
    TIME=`date +"%Y%m%d-%H:%M:%S"` 
    if [ $# -lt 2 ];then
        echo "send email fail,arg:$*"|mail -s "[NOTICE][$MOUDLE_ANME][$TIME][send_email fails]"
    fi  
    SUBJECT=[$TIME]$1
    CONT=$2
    TOS=$3
    for to in $TOS
        do  
            echo $CONT | mail -s "$SUBJECT"  $to 
        done
}

