document.addEventListener("DOMContentLoaded", function(){
    var scrollbar = document.body.clientWidth - window.innerWidth + 'px';
    console.log(scrollbar);
    document.querySelector('[href="#openModal"]').addEventListener('click',function(){
      document.body.style.overflow = 'hidden';
      document.querySelector('#openModal').style.marginLeft = scrollbar;
    });
    document.querySelector('[href="#close"]').addEventListener('click',function(){
      document.body.style.overflow = 'visible';
      document.querySelector('#openModal').style.marginLeft = '0px';
    });
  });

  //Выпадающий список
$("li").mouseover(function(){
    $(this).find('.drop-down').slideDown(300);
    $(this).find(".accent").addClass("animate");
    $(this).find(".item").css("color","#FFF");
   }).mouseleave(function(){
     $(this).find(".drop-down").slideUp(300);
      $(this).find(".accent").removeClass("animate");
      $(this).find(".item").css("color","#000");
   });