import React, { useRef } from 'react'
import { useParams } from 'react-router-dom';
import { bufferDecode, bufferEncode } from '../../utils/helper'

const TotpFallbackPage = () => {
    const { username } = useParams();
    
    function registerUser() {
        if (username === "") {
            alert("Invalid username");
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

    const inputEl = useRef("")
    const onVerify = () => {
        fetch(`http://localhost:3000/otp/verify/${username}`, {
            method: 'POST',
            headers: {
            'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                nonce: inputEl.current.value,
            })
        })
        .then(res => res.json())
        .then((verified) => {
            if(verified){
                registerUser()
            } else {
                alert("Invalid OTP!")
            }
        })
        .catch((err) => {
            console.log("OTP Verification error:" + err)
        })
    }
  return (
    <div>
        <input ref={inputEl} type='text' placeholder='Enter the TOTP'/> 
        <button onClick={onVerify}>Verify</button>
    </div>
  )
}

export default TotpFallbackPage