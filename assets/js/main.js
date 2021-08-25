function cpuTemp(usage){
    sliced=usage.slice(0,-1);
    
    switch(true){
        case sliced <=40:
        return "0,200,10";
        break;

        case sliced >40 && sliced <79:
        return "250,184,0";
        break;

        case sliced >=80:
        return "250,0,0";
        break;

    }
    
    
}