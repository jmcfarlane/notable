/**
 * @fileoverview Describes a bootstrap tab
 */
define([
  'text!templates/noteDetail.html',
  'text!templates/tab.html',
  'backbone',
  'underscore',
  'lib/jquery.hotkeys',
  'codemirror',
  'lib/codemirror/mode/rst/rst'
],
function(noteDetailTemplate, tabTemplate) {
  return Backbone.View.extend({

    initialize: function() {
      $(document).bind('keydown', 'ctrl+s', _.bind(this.save, this));
      this._editor = null;
    },

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

      // Use codemirror for the content
      this.el = $(tabContent).find('.editor').last().parent();
      this.$el = $(this.el);
      this._editor = CodeMirror(_.first(this.$el.find('.editor')), {
        mode: 'rst',
        value: note.content
      });

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

    save: function() {
      if (!this.$el.is(":visible")) {
        return;
      }
      this.model.set('content', this._editor.getValue());
      this.model.set('subject', this.$el.find('.subject').val());
      this.model.set('password', '');
      this.model.save();
      return false;
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
