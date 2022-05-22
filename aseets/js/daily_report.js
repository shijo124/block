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
            let yyyy = d.getFullYear();
            let mm = ( '00' + (d.getMonth() + 1) ).slice( -2 );
            let dd = d.getDate();
            let today = yyyy + '-' + mm + '-' + dd;
            document.getElementById("today").value = today;
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

function input_daily_report(){
    let body = {
        date: document.getElementById("today").value,
        report: document.getElementById("text").value,
    };

    console.log(body);

    fetch("https://dix.api.hello-oi.com/create_daily_report",
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

            var report_end_modal = new bootstrap.Modal(document.getElementById('report_end_modal'), {
                keyboard: false
            })

            document.getElementById("get_coin").textContent = json.mining_coin;

            report_end_modal.show()
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

function get_dix_coin_report(){
    let body = {
    };

    fetch("https://dix.api.hello-oi.com/get_dix_coin_report",
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

            location.href = "https://dix.front.hello-oi.com/block/wallet.html";
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

