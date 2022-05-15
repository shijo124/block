"use strict";

function load_daily_report(){
    let body = {
    };

    fetch("https://dix.api.hello-oi.com/user_wallet",
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
            console.log("if = Trueです");
            console.log(json.message);
            console.log(json.user_name);
            console.log(json.have_coin);

            document.getElementById("user_name").textContent = json.user_name;
            document.getElementById("have_coin").textContent = json.have_coin;

            let d = new Date();
            let today = d.getFullYear() + '-' + d.getMonth() + '-' + d.getDate();
            document.getElementById("have_coin").value = today;
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

