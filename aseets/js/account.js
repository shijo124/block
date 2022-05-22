"use strict";

function create_account(){
    let body = {
        name: document.getElementById("name").value,
        email: document.getElementById("email").value,
        pass: document.getElementById("pass").value,
    };

    fetch("https://dix.api.hello-oi.com/create_account",
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
            location.href = "https://dix.front.hello-oi.com/block/login.html";
        }
        else{
            console.log("if = falseです");
            console.log(json.res_flag);
            console.log(json.message);
        }
    });
}
