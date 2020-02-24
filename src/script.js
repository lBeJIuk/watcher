document.querySelector('#button').addEventListener('click', () => {
  require(['m1/comp1.js', 'm2/comp3.js'], function(comp1, comp3) {
    comp1.log();
    comp3.log();
  });
});