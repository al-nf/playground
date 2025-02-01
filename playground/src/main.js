import { getAuth, signInWithPopup, GoogleAuthProvider } from "firebase/auth";
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

document.getElementById("login").addEventListener("click", () => {
    console.log("button");
    signInWithPopup(auth, provider)
        .then((result) => {
            // This gives you a Google Access Token. You can use it to access the Google API.
            const credential = GoogleAuthProvider.credentialFromResult(result);
            const token = credential.accessToken;
            // The signed-in user info.
            const user = result.user;
            // IdP data available using getAdditionalUserInfo(result)
            console.log(credential,token,user);
            document.getElementById("msg").innerText = "Hi, " + user.displayName;
            fetch()
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
            // ...
        });
    });
