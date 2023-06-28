import React, { useRef } from 'react'
import { bufferDecode, bufferEncode } from '../../utils/helper'
import { Link } from 'react-router-dom'

const Form = () => {
    const emailEl = useRef('')
    //register
    function registerUser() {
        const username = emailEl.current.value
        if (username === "") {
            alert("please enter a username");
            return;
        }

        fetch(`http://localhost:3000/register/begin/${username}`)
        .then(response => response.json())
        .then(async (credentialCreationOptions) => {
            console.log(credentialCreationOptions.publicKey)
            console.log("SERVER'S CHALLENGE----------")
            console.log(credentialCreationOptions.publicKey.challenge)
            credentialCreationOptions.publicKey.attestation = "direct"
            credentialCreationOptions.publicKey.timeout = 300000
            credentialCreationOptions.publicKey.authenticatorSelection = {
                "userVerification": "preferred"
            }
            credentialCreationOptions.publicKey.challenge = bufferDecode(credentialCreationOptions.publicKey.challenge);
            credentialCreationOptions.publicKey.user.id = bufferDecode(credentialCreationOptions.publicKey.user.id);
            console.log(credentialCreationOptions)
            return await navigator.credentials.create({
                publicKey: credentialCreationOptions.publicKey
            })
        })
        .then((credential) => {
            console.log(credential)
            let attestationObject = credential.response.attestationObject;
            let clientDataJSON = credential.response.clientDataJSON;
            const utf8Decoder = new TextDecoder('utf-8');
            const decodedClientData = utf8Decoder.decode(credential.response.clientDataJSON)
            const clientDataObj = JSON.parse(decodedClientData);
            // bufferEncodedClientDataJSON = bufferEncode(clientDataJSON)
            console.log(clientDataObj)
            let rawId = credential.rawId;
            
            const payload = {
                    id: credential.id,
                    rawId: bufferEncode(rawId),
                    type: credential.type,
                    response: {
                        attestationObject: bufferEncode(attestationObject),
                        clientDataJSON: bufferEncode(clientDataJSON),
                    },
                }
            console.log(payload)
            console.log(JSON.stringify(payload))
            fetch(`http://localhost:3000/register/finish/${username}`, {
                method: 'POST',
                headers: {
                'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            })
        })
        .then((success) => {
            alert("successfully registered " + username + "!")
            return
        })
        .catch((error) => {
            console.log(error)
            alert("failed to register " + username)
        })
    }

    //login
    function loginUser() {

        const username = emailEl.current.value
        if (username === "") {
          alert("please enter a username");
          return;
        }

        fetch(`http://localhost:3000/login/begin/${username}`)
        .then(response => {
            if(response.status === 404) {
                throw new Error("User not found");
            }
            if (!response.ok) {
                throw new Error("Failed to login");
            }
            return response.json()
        })
        .then(async (credentialRequestOptions) => {
            credentialRequestOptions.publicKey.challenge = bufferDecode(credentialRequestOptions.publicKey.challenge);
            credentialRequestOptions.publicKey.allowCredentials.forEach(function (listItem) {
                listItem.id = bufferDecode(listItem.id)
            });

            return await navigator.credentials.get({
                publicKey: credentialRequestOptions.publicKey
            })
            .catch((err) => {
                console.log(err)
                fetch(`http://localhost:3000/otp/generate/${username}`)
                .then((res) => res.json())
                .then((data) => {
                    console.log(data.nonce)
                    alert("Your OTP is: " + data.nonce)
                })
                .catch((err) => {
                    console.log("OPT generation failed:" + err)
                })
                alert("failed to login because the device is not registered, first register the device." + username)
                window.location.href = `/totp/${username}`
            })
        })
        .then((assertion) => {
            // TODO
            let authData = assertion.response.authenticatorData;
            let clientDataJSON = assertion.response.clientDataJSON;
            let rawId = assertion.rawId;
            let sig = assertion.response.signature;
            let userHandle = assertion.response.userHandle;
            let response = {
                authenticatorData: bufferEncode(authData),
                clientDataJSON: bufferEncode(clientDataJSON),
                signature: bufferEncode(sig),
                userHandle: bufferEncode(userHandle),
            }
            const payload = {
                id: assertion.id,
                rawId: bufferEncode(rawId),
                type: assertion.type,
                response: response,
            }
            console.log(payload)
            fetch(`http://localhost:3000/login/finish/${username}`, {
                method: 'POST',
                headers: {
                'Content-Type': 'application/json'
                },
                body: JSON.stringify(payload)
            })
            .then((response) => {
                // if(response.status === 401) {

                // }
                if(!response.ok) {
                    throw new Error("Failed to login");
                }
                alert("successfully logged in " + username + "!")
                return
            })
            .catch((error) => {
                console.log(error)
                alert("failed to register " + username)
            })
        })
        .catch((err) => {
            if(err.message === "User not found") {
                alert("User not found. Please register the user.");
            } else {
                console.log(err)
                alert("failed to login " + username)
            }
        })
    }
    
    return (
        <>
            <div>
                Username:
                <br/>
                <input ref={emailEl} type="text" name="username" id="email" placeholder="i.e. foo@bar.com"/>
                <br/>
                <br/>
                <button onClick={() => registerUser()}>Register</button>
                <button onClick={() => loginUser()}>Login</button>
            </div>
        </>
    )
}

export default Form