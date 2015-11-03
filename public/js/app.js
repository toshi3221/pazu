$(function() {
  $('form').ajaxForm({
    success: function (responseJson, statusText) {
       $('#form_result').text(JSON.stringify(responseJson, null, "\t"));
    }
  });
});
