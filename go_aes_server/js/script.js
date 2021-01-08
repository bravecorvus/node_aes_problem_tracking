new Vue ({
      el: '#app',
      data: {
        aes: {
          file: null,
          fileName: '',
        },
      },

      created() {
      },

      methods: {

        storeAESBlob(arg) {
          var file = arg.target.files[0] || arg.dataTranfer.files[0]
          this.aes.file = file;
          this.aes.fileName = file.name;
          document.getElementById('file-input').value = '';
        },

        storeBlowfishBlob(arg) {
          var file = arg.target.files[0] || arg.dataTranfer.files[0]
          this.blowfish.file = file;
          this.blowfish.fileName = file.name;
          document.getElementById('file-input').value = '';
        },

        saveData(data, filename) {
          var a = document.createElement("a");
          document.body.appendChild(a);
          a.style = "display: none";
          blob = new Blob([data], {type: "octet/stream"}),
          url = window.URL.createObjectURL(blob);
          a.href = url;
          a.download = filename;
          a.click();
          window.URL.revokeObjectURL(url);
        },

        aesEncrypt() {
          const formData = new FormData();
          formData.append("file", this.aes.file);
          let itThis = this;
            axios({
              method: 'post',
              url: '/aes_encrypt',
              data: formData,
              headers: {
                'Content-Type': 'multipart/form-data',
              },
              withCredentials: true,
              responseType: 'blob',
            })
            .then(function (result) {
							itThis.saveData(result.data, 'go-encrypted-' + itThis.aes.fileName)
            }, function (error) {
							error.response.data.text().then(text => {
              	Swal.fire({
                	title: 'Error',
                	text: text,
                	type: 'error',
                	confirmButtonText: 'OK'
              	});
							});
            })
        },

        aesDecrypt() {
          const formData = new FormData();
          formData.append("file", this.aes.file);
          let itThis = this;
            axios({
              method: 'post',
              url: '/aes_decrypt',
              data: formData,
              headers: {
                'Content-Type': 'multipart/form-data',
              },
              withCredentials: true,
              responseType: 'blob',
            })
            .then(function (result) {
							itThis.saveData(result.data, 'go-decrypted-' + itThis.aes.fileName)
            }, function (error) {
							error.response.data.text().then(text => {
              	Swal.fire({
                	title: 'Error',
                	text: text,
                	type: 'error',
                	confirmButtonText: 'OK'
              	});
							});
            })
        },

      },


});
