"use strict";

function get_test(){
    console.log(document.getElementById("email").value);
    console.log(document.getElementById("pass").value);
    let body = {
        email: document.getElementById("email").value,
        pass: document.getElementById("pass").value,
    };

    fetch("https://dix.api.hello-oi.com/login",
        {
            headers: {
                "Content-Type": "application/json",
            },
            method: "POST",
            mode: "cors",
            credentials: "include",
            body: JSON.stringify(body),
        }
    )
    .then((res) => {
        return res.json();
    })
    .then((json) => {
        if(json.res_flag === true ){
            location.href = "https://dix.front.hello-oi.com/block/daily_report.html";
        }
        else{
            console.log("if = falseです");
            console.log(json.res_flag);
            console.log(json.message);
            console.log(json.user.Id);
            console.log(json.user.Name);
            console.log(json.user.Email);
            console.log(json.user.Pass);
        }
    });
}

