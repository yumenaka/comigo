<template>
    <div>
        <!-- @submit.prevent   .prevent 表示提交以后不刷新页面，prevent是preventDefault,阻止标签默认行为，有些标签有默认行为，例如a标签的跳转链接属性href等。 -->
        <form @submit.prevent="submit">
            <input type="text" v-model="username" placeholder="Username" required />
            <input type="password" v-model="password"  volarvolar placeholder="Password" required />
            <button type="submit">Login</button>
        </form>
    </div>
</template>

<script lang="ts">
import axios from "axios";

export default {
    data() {
        return {
            username: "",
            password: "",
        };
    },
    methods: {
        async submit() {
            try {
                const response = await axios.post("/api/login", {
                    username: this.username,
                    password: this.password,
                });
                console.log(response.data);
                // 做登录成功后的操作,例如跳转页面
            } catch (error) {
                // 做登录失败时的操作
                console.log(error);
            }
        },
    },
};
</script>
