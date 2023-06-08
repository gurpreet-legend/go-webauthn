import React from 'react'

const Form = () => {
    const registerUser = () => {

    }

    const loginUser = () => {

    }
    
    return (
        <div>
            Username:
            <br/>
            <input type="text" name="username" id="email" placeholder="i.e. foo@bar.com"/>
            <br/>
            <br/>
            <button onClick={registerUser()}>Register</button>
            <button onClick={loginUser()}>Login</button>
        </div>
    )
}

export default Form