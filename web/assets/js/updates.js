// Updates css colors if a server status changes, should be called from main.js each time servers 
// are pinged for status


function update_css(){
    var status = document.getElementsByClassName("status");
    

    let text = "";
    for (let x in status) {
    text = status.innerHTML;
    console.log(text);
    }
}