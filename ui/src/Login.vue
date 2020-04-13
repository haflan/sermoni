<template>
    <div style="width: 100%;">
        <div style="max-width: 400px;">
            <input ref="ppinput" 
                   v-model="passphrase"
                   type="password"
                   placeholder="passphrase">
            <button @click="enter">&gt;</button>
        </div>
    </div>
</template>

<script>
    import api from "./requests.js";
    export default {
        name: "Login",
        data() {
            return {
                passphrase: ""
            }
        },
        methods: {
            enter() {
                api.login(
                    this.passphrase,
                    success => {
                        this.$emit("login");
                    },
                    error => {
                        this.$emit("error");
                    }
                );
            }
        },
        mounted() {
            this.$refs.ppinput.focus();
            this.$refs.ppinput.addEventListener("keypress", (e) => {
                if (e.keyCode === 13) {
                    e.preventDefault();
                    this.enter();
                }
            });
        }
    }

</script>

<style scoped>
div {
    display: flex;
    align-items: center
}
input {
    border-width: 0 0 1px;
    border-color: #bbf;
    background-color: inherit;
    outline-width: 0;
    flex: 1;
    margin-bottom: 3px;
    font-size: 2em;
    width: 100%;
    max-width: 400px;
}
input:focus {
    outline-width: 0;
}
input::placeholder {
    color: #bbf;
}
button {
    background-color: inherit;
    border: none;
    border-radius: 3px;
    color: #bbf;
    outline-width: 0;
    text-decoration: none;
    text-align: center;
}
</style>
