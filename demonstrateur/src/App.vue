<template>
  <div class="container">
    <div class="wrapper">
      <h1>Download file</h1>
      <div class="item">
        <form @submit.prevent="downloadFile">
          <input
            type="text"
            id="filename"
            placeholder="Your file.."
            v-model="filename"
            required
          />
          <button type="submit">Download</button>
        </form>
      </div>
    </div>
    <div class="wrapper">
      <h1>Upload file</h1>

      <div class="item">
        <form @submit.prevent="uploadFile">
          <input type="file" id="file" ref="fileInput" required />
          <div class="input-upload">Select your file</div>
          <button type="submit">Upload it</button>
        </form>
      </div>
    </div>
  </div>
</template>
<style>
body {
  margin: 0;
  padding: 0;
  background: rgb(6, 1, 54);
  font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
}

form {
  position: relative;
}
input[type="file"] {
  opacity: 0;
}

textarea:focus, input:focus{
  outline: none;
}

.input-upload {
  position: absolute;
  display: flex;
  justify-content: center;
  align-items: center;
  color: black;
  left: 20%;
  top: 0;
  background: rgb(255, 243, 243);
  border-radius: 5px;
  pointer-events: none;
  width: 45%;
  height: 100%;
  overflow: hidden;
}


h1 {
  background: linear-gradient(to right, #0fbbff, #7e0fff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.container .wrapper {
  width: 35%;
  text-align: center;
}
.container .item:after {
  position: absolute;
  content: "";
  top: 5vw;
  left: 0;
  right: 0;
  z-index: -1;
  height: 100%;
  width: 100%;
  margin: 0 auto;
}

.container {
  padding: 0 20px;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 30px;
}

.container .wrapper .item {
  position: relative;
  text-align: center;
  background: rgb(51, 33, 212);
  animation: textColor 10s ease infinite;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 300px;
  padding: 10px;
  border-radius: 10px;
}
.container .wrapper .item form {
  display: flex;
  justify-content: center;
  align-items: stretch;
}
.container .item label {
  margin-right: 10px;
}
.container .item input {
  padding: 10px 0 10px 20px;
  font-size: 16px;
  border: none;
  border-radius: 10px;
  margin-right: 10px;
}
.container .item button {
  padding: 10px 20px;
  border-radius: 5px;
  border: none;
  background: white;
  transition: all .3s ease;
}
.container .item button:hover {
  transform: scale(1.1);
}
</style>
<script>
import { ref } from "vue";
import { useToast } from "vue-toastification";

export default {
  setup() {
    const filename = ref("");
    const fileInput = ref(null);
    const toast = useToast();

    // Fonction pour télécharger un fichier
    async function downloadFile() {
      try {
        console.log(filename);
        const response = await fetch(
          `http://localhost:8080/download/${filename.value}`
        );
        if (response.status === 200) {
          const blob = await response.blob();
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = filename;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);
          toast.success("File downloaded successfully");
        } else {
          toast.error(
            `Impossible to download the file "${filename.value}"`
          );
        }
      } catch (error) {
        console.error(error);
      }
    }

    // Fonction pour téléverser un fichier
    async function uploadFile() {
      const file = fileInput.value.files[0];
      const formData = new FormData();
      formData.append("file", file);

      try {
        const response = await fetch("http://localhost:8080/upload", {
          method: "POST",
          body: formData,
        });
        const data = await response.json();
        toast.success("File uploaded successfully");
      } catch (error) {
        toast.error("Something went wrong");
        console.error(error);
      }
    }

    return { filename, fileInput, downloadFile, uploadFile };
  },
};
</script>
