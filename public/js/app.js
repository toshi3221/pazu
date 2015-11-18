$(function() {
  $('form').ajaxForm({
    success: function (response, statusText, xhr, element) {
       $('#form_result').html(xhr.responseText);
    }
  });
});
