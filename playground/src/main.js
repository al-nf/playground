import { getAuth, signInWithPopup, GoogleAuthProvider, onAuthStateChanged } from "firebase/auth";
import { initializeApp } from 'firebase/app';

const firebaseConfig = {
  projectId: "test-go-457f8",
  apiKey: "AIzaSyA8COMPO956Uze1ncOAPc84srV11Crk8F0",
  authDomain: "test-go-457f8.firebaseapp.com",
  storageBucket: "test-go-457f8.firebasestorage.app",
  messagingSenderId: "60184479922",
  appId: "1:60184479922:web:53114e2df89f49ca7121ab"
};

const app = initializeApp(firebaseConfig);

const auth = getAuth();
const provider = new GoogleAuthProvider();


const button = document.getElementById("login");

onAuthStateChanged(auth, (user) => {
  if (user) {
    document.getElementById("msg").innerText = "Hi, " + user.displayName;
    button.innerText = "Log Out";
  }
  else {
    document.getElementById("msg").innerText = "Goodbye";
    button.innerText = "Log In";
  }
})

button.addEventListener("click", (e) => {
    e.preventDefault();
    if(auth.currentUser){
      auth.signOut();
    }
    else {
      login();
    }
});


function login() {
  console.log("button");
    signInWithPopup(auth, provider)
        .then(async (result) => {
            // This gives you a Google Access Token. You can use it to access the Google API.
            const credential = GoogleAuthProvider.credentialFromResult(result);
            const token = credential.accessToken;
            // The signed-in user info.
            const user = result.user;
            // IdP data available using getAdditionalUserInfo(result)
            console.log(credential,token,user);
            
            const idToken = await user.getIdToken();
            sendToken(idToken);
        }).catch((error) => {
            console.log("oh naur");
            // Handle Errors here.
            const errorCode = error.code;
            const errorMessage = error.message;
            console.log(errorCode,errorMessage);
            // The email of the user's account used.
            const email = error.customData.email;
            // The AuthCredential type that was used.
            const credential = GoogleAuthProvider.credentialFromError(error);
        });
}
async function sendToken(token) {
  var myHeaders = new Headers();
  myHeaders.append("Authorization", "Bearer "+token);

  var requestOptions = {
    method: 'POST',
    headers: myHeaders,
    redirect: 'follow'
  };
  const response = await fetch("http://localhost:8080/auth", requestOptions)
    .then(response => response.text())
    .then(result => console.log(result))
    .catch(error => console.log('error', error));
}
