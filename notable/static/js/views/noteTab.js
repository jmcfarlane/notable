/**
 * @fileoverview Describes a bootstrap tab
 */
define([
  'text!templates/noteDetail.html',
  'text!templates/tab.html',
  'backbone',
  'underscore'
],
function(noteDetailTemplate, tabTemplate) {
  return Backbone.View.extend({

    render: function(collection) {
      var note = this.model.toJSON(),
        tabs = this.options.tabs,
        tabContent = this.options.tabContent;

      // No tabs are active
      tabs.find('.active').removeClass('active');

      // Add new tab content
      tabContent.append(_.template(noteDetailTemplate, {
        note: note
      }));

      // Add a new tab
      tabs.append(_.template(tabTemplate, {
        note: note
      }));

      // Activate the new tab
      tabs.find(this.selector())
        .tab('show')
        .find('button')
        .on('click', _.bind(this.close, this));

      return this;
    },

    selector: function() {
      return 'a[href=#'+ this.model.get('uid') +']';
    },

    close: function() {
      // Show the previous tab and then remove this one
      this.options.tabs.find(this.selector())
        .parent()
        .prev()
        .find('a')
        .tab('show')
        .parent()
        .next()
        .remove();

      // Remove the tab content (after the ^ animation is finished) (ick)
      setTimeout(_.bind(function() {
        $('#' + this.model.get('uid')).detach().remove();
      }, this), 300);
    }

  });
});
