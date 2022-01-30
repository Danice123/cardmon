import React from 'react'
import ReactDOM from 'react-dom'

function App() {
    return (
        <h1>Hello world</h1>
    )
}

const root = document.createElement("div")
document.getElementsByTagName('body')[0].prepend(root)
ReactDOM.render(<App />, root)